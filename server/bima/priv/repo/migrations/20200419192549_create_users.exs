defmodule Bima.Repo.Migrations.CreateUsers do
  use Ecto.Migration

  def change do
    create table(:users) do
      add :email, :string
      add :encrypted_password, :string
      add :salt, :string

      timestamps()
    end

    create unique_index(:users, [:email])
  end
end
