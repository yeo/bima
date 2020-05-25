defmodule BimaWeb.Api.BlobView do
  use BimaWeb, :view

  def render("error.json", %{error: error}) do
    %{
      "error": error
    }
  end
end
