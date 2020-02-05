defmodule Bima.Token do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, Ecto.UUID, autogenerate: false}
  schema "tokens" do
    field :name, :string
    field :token, :string
    field :url, :string
    field :version, :integer

    field :app_id, :string
    field :deleted_at, :time

    timestamps()
  end

  @doc false
  def changeset(token, attrs) do
    token
    |> cast(attrs, [:id, :name, :url, :token, :version, :app_id])
    |> validate_required([:id, :name, :url, :token, :version, :app_id])
  end

end
