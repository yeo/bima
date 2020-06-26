defmodule BimaWeb.Api.MeController do
  use BimaWeb, :controller

  import Ecto.Query
  alias Bima.{Repo}
  alias Bima.Apps

  alias Plug.Conn

  def index(conn, _params) do
    [app_id | _ ] = Conn.get_req_header(conn, "appid")
    [app_version | _ ] = Conn.get_req_header(conn, "appversion")

    [major, minor, patch] = String.split(app_version, ".")

    IO.inspect app_id, label: 'app_id'
    IO.inspect app_version, label: 'app_version'

    #changes = App.changeset(%App{}, %{id: app_id})
    app = Apps.ensure_app(app_id)
    IO.inspect app, label: 'app'
    Apps.bump_db(app_id)

    resp = %{
      news: [%{body: "this is a test", url: "here", id: 1}],
      db: app.db_version,
    }

    if true do
      resp = Map.put(resp, :update, %{version: "1.2", url: "go here to download"})
    end

    conn
    |> json(resp)
  end
end
