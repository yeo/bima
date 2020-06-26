defmodule Bima.Repo.Migrations.CreateNews do
  use Ecto.Migration

  def change do
    create table(:news) do
      add :body, :text

      timestamps()
    end

  end
end
