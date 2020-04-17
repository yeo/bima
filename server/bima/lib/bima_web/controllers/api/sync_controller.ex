defmodule BimaWeb.Api.SyncController do
  use BimaWeb, :controller

  import Ecto.Query
  alias Bima.{Repo, Token}

  alias Plug.Conn

  def sync(conn, %{"current" => client_tokens, "removed" => request_removed_tokens}) do
    [app_id | _ ] = Conn.get_req_header(conn, "appid")
    tokens = Repo.all(from u in Token, where: u.app_id == ^app_id)
    sync_add_or_update(tokens, client_tokens, app_id)
    removed_tokens = sync_remove(client_tokens, request_removed_tokens)

    tokens = Repo.all(from u in Token, where: u.app_id == ^app_id and is_nil(u.deleted_at))
    render(conn, "sync.json", tokens: tokens, removed_tokens: removed_tokens, added: added_tokens)
  end


  defp sync_add_or_update(exist_tokens, client_tokens, app_id) do
    Enum.each(client_tokens, fn x -> sync_token(exist_tokens, x, app_id) end)
  end

  defp sync_add(exist_tokens) do
  end

  defp sync_update(exist_tokens) do
  end

  defp sync_remove(current_client_tokens, request_removed_tokens) do
    removed_tokens = Enum.map(request_removed_tokens, fn request_token ->
      token = Repo.get(Token, request_token["id"])
      if token do
        r = Ecto.Changeset.cast(token, %{deleted_at: DateTime.utc_now(), version: token.version + 1}, [:deleted_at, :version])
        |> Repo.update!()
        IO.inspect r

        %{id: token.id, version: token.version}
      else
        %{id: request_token["id"], version: -1}
      end
    end)

    no_longer_exist_token = Repo.all(from u in Token, select: [:id, :version], where: not(is_nil(u.deleted_at)) and u.id in ^Enum.map(current_client_tokens, &(&1["id"])))
    removed_tokens ++ Enum.map(no_longer_exist_token, &(%{id: &1.id, version: &1.version}))
  end

  defp sync_token(exist_tokens, new_token, app_id) do
    exist_token = Enum.at(Enum.filter(exist_tokens, fn x -> new_token["id"] == x.id end), 0)

    if exist_token == nil do
        # These are new token that client submit and does't existed in our db so we add them in
        changeset = Token.changeset(%Token{app_id: app_id}, new_token)
        Repo.insert(changeset)
    else
        IO.puts "exist token"
        IO.inspect exist_token
        # Client version is higher than server, which mean it's changed on client side
        # TODO: handle case of same version mean 2 client has update, we can also check updated time
        if exist_token.deleted_at do
          IO.puts "Token #{exist_token.id} is request but it's mark for deletion at #{exist_token.deleted_at}"
        else
          if new_token["version"] > exist_token.version do
            changeset = Ecto.Changeset.change(exist_token,
              name: new_token["name"],
              url: new_token["url"],
              version: new_token["version"])
            Repo.update changeset
          end
        end
    end
  end
end
