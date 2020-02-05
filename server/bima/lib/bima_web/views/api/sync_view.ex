defmodule BimaWeb.Api.SyncView do
  use BimaWeb, :view

  def render("sync.json", %{tokens: tokens}) do
    %{
      data: Enum.map(tokens, &token_json/1)
    }
  end

  def token_json(token) do
    %{
      id: token.id,
      name: token.name,
      url: token.url,
      token: token.token,
      version: token.version,
    }
  end
end
