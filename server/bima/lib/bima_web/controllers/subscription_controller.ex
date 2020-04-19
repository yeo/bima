defmodule BimaWeb.SubscriptionController do
  use BimaWeb, :controller

  def index(conn, _params) do
    render(conn, "index.html")
  end
end
