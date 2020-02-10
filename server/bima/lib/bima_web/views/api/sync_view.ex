defmodule BimaWeb.Api.SyncView do
  use BimaWeb, :view

  def render("sync.json", %{tokens: tokens, removed_tokens: removed_tokens}) do
    %{
      current: Enum.map(tokens, &token_json/1),
      removed: removed_tokens,
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
