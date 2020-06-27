defmodule BimaWeb.Plugs.Guest do
  import Plug.Conn
  import Phoenix.Controller

  def init(opts), do: opts

  def call(conn, _opts) do
    if Plug.Conn.get_session(conn, :current_user_id) do
      conn
      |> put_flash(:info, "You already login")
      |> redirect(to: BimaWeb.Router.Helpers.page_path(conn, :index))
      |> halt()
    end
    conn
  end
end
