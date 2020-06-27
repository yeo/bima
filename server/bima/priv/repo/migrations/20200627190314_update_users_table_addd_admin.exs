defmodule Bima.Repo.Migrations.UpdateUsersTableAdddAdmin do
  use Ecto.Migration

  def change do
    alter table("users") do
      add :admin, :boolean, defaul: false
    end
  end
end
