defmodule Bima.Repo.Migrations.CreateApps do
  use Ecto.Migration

  def change do
    create table(:apps, primary_key: false) do
      add :id, :uuid, primary_key: true
      add :db_version, :integer, default: 1
      add :db_updated_at, :utc_datetime, null: true

      timestamps()
    end

  end
end
