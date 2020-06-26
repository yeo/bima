defmodule Bima.Client.News do
  use Ecto.Schema
  import Ecto.Changeset

  schema "news" do
    field :body, :string

    timestamps()
  end

  @doc false
  def changeset(news, attrs) do
    news
    |> cast(attrs, [:id, :body])
    |> validate_required([:id, :body])
  end
end
