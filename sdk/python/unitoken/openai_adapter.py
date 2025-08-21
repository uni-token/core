"""OpenAI adapter for UniToken service."""

from .client import UniTokenClient


def request_unitoken_openai(
    app_name: str, description: str, saved_api_key: str | None = None
) -> tuple[str, str | None]:
    """
    Requests user for OpenAI token via UniToken service.
    
    This function registers an application with the UniToken service and retrieves
    an OpenAI-compatible API endpoint and token. If the user has not granted permission,
    the API key will be None but the base URL will still be provided.
    
    Args:
        app_name (str): The name of the application requesting access.
        description (str): A description of what the application does.
        saved_api_key (str | None, optional): Previously saved API key for the app.
            If provided, it will be used to identify the existing app registration.
            Defaults to None.
    
    Returns:
        tuple[str, str | None]: A tuple containing:
            - base_url (str): The OpenAI-compatible API endpoint URL
            - api_key (str | None): The API token for authentication, or None if
              the user has not granted permission
    
    Raises:
        Exception: If the UniToken service fails to start, network issues occur,
            or service errors are encountered during app registration.
    """
    client = UniTokenClient()
    base_url = str(client.server_url) + "openai/"

    response = client.post(
        "app/register",
        json={
            "name": app_name,
            "description": description,
            "uid": saved_api_key,
        },
        timeout=None,
    )

    if response.status_code == 403:
        return base_url, None

    if not response.ok:
        raise Exception(
            f"Failed to register app: {response.status_code} {response.text}"
        )

    api_key = response.json().get("token", None)
    return base_url, api_key
