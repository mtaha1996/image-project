from rest_framework import permissions
from django.conf import settings
from .models import Access


def is_valid_api_key_and_engin_id(request):
    api_key = request.query_params.get("key")
    engin_id = request.query_params.get("cx")

    if api_key and engin_id:
        return Access.objects.filter(api_key=api_key, engin_id=engin_id).exists()
    return False


class IsValidApiKeyAndEnginId(permissions.BasePermission):
    def has_permission(self, request, view):
        return is_valid_api_key_and_engin_id(request)
