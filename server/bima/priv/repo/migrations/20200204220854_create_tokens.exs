defmodule Bima.Repo.Migrations.CreateTokens do
  use Ecto.Migration

  @disable_migration_lock true

  def change do
    create table(:tokens, primary_key: false) do
      add :id, :uuid, primary_key: true
      add :name, :string
      add :url, :string
      add :token, :string
      add :version, :integer

      add :app_id, :string
      add :deleted_at, :time

      timestamps()
    end

    create index("tokens", [:app_id], concurrently: true)
    create index("tokens", [:deleted_at], concurrently: true)
  end
end
