defmodule Bima.Repo.Migrations.CreateBlobs do
  use Ecto.Migration

  def change do
    create table(:blobs) do
      add :code, :string
      add :payload, :string
      add :ttl, :integer

      timestamps()
    end

    create unique_index("blobs", [:code])
  end
end
