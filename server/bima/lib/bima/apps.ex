defmodule Bima.Apps do
  @moduledoc """
  The Accounts context.
  """

  import Ecto.Query, warn: false
  alias Bima.Repo

  alias Bima.Client.App

  def list_apps do
    Repo.all(App)
  end

  def get_app!(id), do: Repo.get!(App, id)

  def create_app(attrs \\ %{}) do
    %App{}
    |> App.changeset(attrs)
    |> Repo.insert()
  end

  def ensure_app(id) do
    case Repo.get(App, id) do
      nil ->
        %App{}
        |> App.changeset(%{id: id})
        |> Repo.insert!(
             on_conflict: :nothing,
             conflict_target: :id
           )
      app -> app
    end
  end

  @doc """
  Bump DB version for an app

  DB version is how we tell client that it needs to fetch new change from server
  """
  def bump_db(app_id) do
    from(a in App, where: a.id == ^app_id, update: [inc: [db_version: 1], set: [db_updated_at: fragment("NOW()")]])
    |> Repo.update_all([])
  end

end
