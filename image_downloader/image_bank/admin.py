from django.contrib import admin
from .models import Image, Access


# Register your models here.
class ImageAdmin(admin.ModelAdmin):
    list_display = ("id", "query", "link")
    search_fields = ("query", "link")
    list_filter = ("query", "link")


class AccessAdmin(admin.ModelAdmin):
    list_display = ("engin_id", "api_key")
    search_fields = ("engin_id", "api_key")
    list_filter = ("engin_id", "api_key")


admin.site.register(Image, ImageAdmin)
admin.site.register(Access, AccessAdmin)
