defmodule Bima.Repo.Migrations.CreateTokens do
  use Ecto.Migration

  @disable_migration_lock true
  @disable_ddl_transaction true

  def change do
    create table(:tokens, primary_key: false) do
      add :id, :uuid, primary_key: true
      add :name, :string
      add :url, :string
      add :token, :string
      add :version, :integer

      add :app_id, :string
      add :deleted_at, :utc_datetime

      timestamps()
    end

    create index("tokens", [:app_id, :deleted_at], concurrently: true)
  end
end
