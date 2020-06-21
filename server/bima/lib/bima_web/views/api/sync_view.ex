defmodule BimaWeb.Api.SyncView do
  use BimaWeb, :view

  def render("sync.json", %{removed: removed_tokens, added: added_tokens, changed: changed_tokens}) do
    %{
      added: Enum.map(added_tokens, &token_json/1),
      changed: Enum.map(changed_tokens, &token_json/1),
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
