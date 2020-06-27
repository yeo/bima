defmodule BimaWeb.Plugs.Auth do
  import Plug.Conn
  import Phoenix.Controller

  alias Bima.Accounts

  def init(opts), do: opts

  def call(conn, opts) do
    if user_id = Plug.Conn.get_session(conn, :current_user_id) do
      current_user = Accounts.get_user!(user_id)
      if !opts[:admin] || (opts[:admin] && current_user.admin) do
        conn |> assign(:current_user, current_user)
      else
        conn
        |> put_flash(:info, "You are not an admin.")
        |> redirect(to: "/")
        |> halt()
      end
    else
      conn
      |> redirect(to: "/login")
      |> halt()
    end
  end
end
