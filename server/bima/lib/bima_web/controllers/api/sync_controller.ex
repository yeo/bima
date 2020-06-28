defmodule BimaWeb.Api.SyncController do
  use BimaWeb, :controller

  import Ecto.Query
  alias Bima.{Repo, Token, Apps}

  alias Plug.Conn

  def sync(conn, %{"current" => client_submit_tokens, "removed" => request_removed_tokens}) do
    [app_id | _ ] = Conn.get_req_header(conn, "appid")
    IO.inspect app_id, label: 'app_id'

    # Apply the removed token to our db
    # This need to be done first to avoid client revert each others
    removed_tokens = sync_remove(client_submit_tokens, request_removed_tokens, app_id)

    # Refresh state to see what we have on our db right now
    current_tokens_map = token_list_as_map(app_id)

    conn
    |> render("sync.json",
      removed: removed_tokens,
      added:   sync_add(current_tokens_map, client_submit_tokens, app_id),
      changed: sync_update(current_tokens_map, client_submit_tokens, app_id))
  end

  # Return a map of tokens where the map key is token id
  #
  # This helps speed up token existed check to O(1)
  defp token_list_as_map(app_id) do
    Repo.all(from u in Token, where: u.app_id == ^app_id and is_nil(u.deleted_at))
    |> Map.new(fn token -> {token.id, token} end)
  end

  defp deleted_token_list_as_map(app_id) do
    Repo.all(from u in Token, where: u.app_id == ^app_id and not is_nil(u.deleted_at))
    |> Map.new(fn token -> {token.id, token} end)
  end

  # Given two list of:
  # - current_client_tokens: the state of local db from client site
  # - current_tokens: the state of remote db from server side
  #
  # This does 2 things:
  # - write the token that existed on client but not on server
  # return list of token we need to add to client: these are tokens existed on server but not from client
  defp sync_update(exist_tokens_map, client_tokens, app_id) do
    client_tokens
    |> Enum.filter(&(Map.has_key?(exist_tokens_map, &1["id"])))
    |> Enum.map(fn new_token ->
      # Client version is higher than server, which mean it's changed on client side
      # TODO: handle case of same version mean 2 client has update, we can also check updated time
      # client version is higher than server state, so we update server state. In this case, no need to return to client, it has the lastest state
      exist_token = Map.get(exist_tokens_map, new_token["id"])
      cond do
        exist_token.version < new_token["version"] ->
          changeset = exist_token
          |> Ecto.Changeset.change(name: new_token["name"], url: new_token["url"], version: new_token["version"])
          case Repo.update(changeset) do
            {:ok, token} ->
              Apps.bump_db(app_id)
              %{token | token: nil}
          end

        exist_token.version > new_token["version"] ->
          %{exist_token | token: nil}

        true ->
          # TODO: handle this. do nothing for now
          nil
      end
    end)
    |> Enum.filter(&(&1))
  end

  # Do 2 things:
  # - Add tokens that don't exist in our database
  # - Return token that exist in our db but not in client request
  defp sync_add(exist_tokens_map, client_tokens, app_id) do
    removed_tokens = deleted_token_list_as_map(app_id)

    client_tokens
    |> Enum.filter(fn client_submited_token -> !Map.get(exist_tokens_map, client_submited_token["id"]) end)
    |> Enum.each(fn new_token ->
      # These are new token that client submit and does't existed in our db so we add them in
      # but we don't need to return them
      check_removed_token = Map.get(removed_tokens, new_token["id"])
      if !check_removed_token do
        changeset = Token.changeset(%Token{app_id: app_id}, new_token)
        case Repo.insert(changeset) do
          {:ok, token} ->
            Apps.bump_db(app_id)
            token
        end
      end
    end)

    client_tokens_map = Map.new(client_tokens, fn token -> {token["id"], token} end)

    IO.inspect(client_tokens_map)

    Enum.map(exist_tokens_map, fn({id, existed_token}) -> existed_token end)
    |> Enum.filter(fn t -> !Map.has_key?(client_tokens_map, t.id) end)
  end

  # Given a list of token from local db of a client, and a list of removed token
  # Return a list of token that is removed, either from client request or removed by other clients that has same app id
  defp sync_remove(current_client_tokens, request_removed_tokens, app_id) do
    removed_tokens = Enum.map(request_removed_tokens, fn request_token ->
      token = Repo.get(Token, request_token["id"])
      if token do
        r =
          Ecto.Changeset.cast(token, %{deleted_at: DateTime.utc_now(), version: token.version + 1}, [:deleted_at, :version])
          |> Repo.update!()

        if r do
          Apps.bump_db(app_id)
        end
        %{id: token.id, version: token.version}
      else
        %{id: request_token["id"], version: -1}
      end
    end)

    # Let's say client A and B share same app ID.
    # client A deleted T1 and stay off line for a while
    # During that, client B deleted T2.
    # Later on, client A come up we will need to tell client A that T2 is deleted too
    no_longer_exist_token = Repo.all(from u in Token, select: [:id, :version], where: not(is_nil(u.deleted_at)) and u.id in ^Enum.map(current_client_tokens, &(&1["id"])))
    removed_tokens ++ Enum.map(no_longer_exist_token, &(%{id: &1.id, version: &1.version}))
  end
end
