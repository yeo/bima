defmodule BimaWeb.Api.MeController do
  use BimaWeb, :controller

  import Ecto.Query
  alias Bima.{Repo}
  alias Bima.Apps

  alias Plug.Conn

  def index(conn, _params) do
    [app_id | _ ] = Conn.get_req_header(conn, "appid")
    [app_version | _ ] = Conn.get_req_header(conn, "appversion")

    IO.inspect app_id, label: 'app_id'
    IO.inspect app_version, label: 'app_version'

    [major, minor, patch] = String.split(app_version, ".") |> Enum.map(&String.to_integer/1)

    #changes = App.changeset(%App{}, %{id: app_id})
    app = Apps.ensure_app(app_id)
    IO.inspect app, label: 'app'

    resp = %{
      news: [%{body: "this is a test", url: "here", id: 1}],
      db: app.db_version,
    }

    resp = if check_version(major, minor, patch), do: Map.put(resp, :update, %{version: "1.2", url: "go here to download"}), else: resp

    conn |> json(resp)
  end

  def check_version(major, minor, patch) do
    [latest: [current_major, current_minor, current_patch], url: _] = Application.get_env(:bima, :bima_version)
    cond do
      current_major > major ->
        true
      current_major < major ->
        false
      current_minor > minor ->
        true
      current_minor < minor ->
        false
      current_patch > patch ->
        true
      current_patch < patch ->
        false
      true ->
        false
    end
  end
end
