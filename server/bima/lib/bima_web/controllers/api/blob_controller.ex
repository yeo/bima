defmodule BimaWeb.Api.BlobController do
  use BimaWeb, :controller

  import Ecto.Query
  alias Bima.{Repo, Blob}

  alias Plug.Conn

  @char_dicts {"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

  def create(conn, %{"payload" => payload}) do
    IO.inspect payload
    case gen_code(payload) do
      {:ok, blob} ->
        json conn, %{"payload" => payload, "code" => blob.code}
      {:error, changset} ->
        json conn, %{"error" => "Fail to create blob"}
    end
  end

  def show(conn, %{"code" => code}) do
    blob = Repo.get_by!(Blob, code: code)

    if NaiveDateTime.diff( NaiveDateTime.utc_now(), blob.inserted_at) < blob.ttl do
      json conn, %{"payload": blob.payload}
    else
      render conn, "error.json", %{error: "Invalid code"}
    end
  end


  defp gen_code(payload) do
    code = for i <- 0..12, i > 0, do: elem(@char_dicts, :rand.uniform(34))
    %Blob{}
      |> Blob.changeset(%{"payload" => payload, "ttl" => 6000, "code" => Enum.join(code, "")})
      |> Repo.insert()
  end
end
