defmodule Bima.News do
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

  def list_news(attrs \\ %{}) do
    %App{}
    |> App.changeset(attrs)
    |> Repo.insert()
  end

end
