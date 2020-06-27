defmodule Bima.Accounts.User do
  use Ecto.Schema
  import Ecto.Changeset

  alias Bima.Accounts.{User, Encryption}

  schema "users" do
    field :email, :string
    field :encrypted_password, :string
    field :salt, :string
    field :admin, :boolean

    # VIRTUAL FIELDS
    field :password, :string, virtual: true
    field :password_confirmation, :string, virtual: true

    timestamps()
  end

  @doc false
  def changeset(user, attrs) do
    user
    |> cast(attrs, [:email, :password])
    |> validate_required([:email])
    |> validate_length(:password, min: 6)
    |> validate_confirmation(:password)
    #|> validate_format(:username, ~r/^[a-z0-9][a-z0-9]+[a-z0-9]$/i)
    |> unique_constraint(:email)
    |> downcase_email
    |> encrypt_password
  end

  defp encrypt_password(changeset) do
    password = get_change(changeset, :password)
    if password do
      encrypted_password = Encryption.hash_password(password)
      put_change(changeset, :encrypted_password, encrypted_password)
    else
      changeset
    end
  end

  defp downcase_email(changeset) do
    update_change(changeset, :email, &String.downcase/1)
  end

  def to_admin_changeset(user) do
    user
    |> cast(%{admin: true}, [:admin])
  end
end
