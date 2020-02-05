defmodule BimaWeb.Api.SyncController do
  use BimaWeb, :controller

  import Ecto.Query
  alias Bima.{Repo, Token}

  alias Plug.Conn

  def sync(conn, %{"current" => client_tokens, "removed" => removed_tokens}) do
    [app_id | tail ] = Conn.get_req_header(conn, "appid")
    tokens = Repo.all(from u in Token, where: u.app_id == ^app_id)

    sync_add_or_update(tokens, client_tokens, app_id)
    sync_remove(removed_tokens)

    tokens = Repo.all(from u in Token, where: u.app_id == ^app_id)
    render(conn, "sync.json", tokens: tokens)
  end


  defp sync_add_or_update(exist_tokens, client_tokens, app_id) do
    Enum.each(client_tokens, fn x -> sync_token(exist_tokens, x, app_id) end)
  end

  defp sync_remove(removed_tokens) do
    for remove_id <- removed_tokens do
      token = Repo.get!(Token, remove_id)
      Repo.delete(token)
    end
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
        if  new_token["version"] > exist_token.version do
          changeset = Ecto.Changeset.change(exist_token,
            name: new_token["name"],
            url: new_token["url"],
            version: new_token["version"])
          Repo.update changeset
        end
    end
  end
end
