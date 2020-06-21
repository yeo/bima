defmodule BimaWeb.Api.FaviconController do
  use BimaWeb, :controller

  import Ecto.Query
  alias Bima.{Repo, Token}

  alias Plug.Conn

  def index(conn, %{"current" => client_tokens, "removed" => request_removed_tokens}) do
  end
end
