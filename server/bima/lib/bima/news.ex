defmodule Bima.News do
  @moduledoc """
  The News context.
  """

  import Ecto.Query, warn: false
  alias Bima.Repo

  alias Bima.Client.News

  def list_news do
    Repo.all(News)
  end

  def get_news!(id), do: Repo.get!(News, id)
end
