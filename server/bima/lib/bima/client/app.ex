defmodule Bima.Client.App do
  use Ecto.Schema
  import Ecto.Changeset
  import Ecto.Query

  @primary_key {:id, Ecto.UUID, autogenerate: true}
  schema "apps" do
    field :db_version, :integer, default: 1
		field :db_updated_at, :utc_datetime
		# field :id, Ecto.UUID

    timestamps()
  end

  @doc false
  def changeset(app, attrs) do
    app
    |> cast(attrs, [:id, :db_version, :db_updated_at])
    |> validate_required([:id])
  end
end
