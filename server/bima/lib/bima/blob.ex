defmodule Bima.Blob do
  use Ecto.Schema
  import Ecto.Changeset

  schema "blobs" do
    field :code, :string
    field :payload, :string
    field :ttl, :integer

    timestamps()
  end

  @doc false
  def changeset(blob, attrs) do
    blob
    |> cast(attrs, [:code, :payload, :ttl])
    |> validate_required([:code, :payload, :ttl])
    |> unique_constraint(:code)
  end
end
