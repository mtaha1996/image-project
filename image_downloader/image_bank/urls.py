from django.urls import path
from .views import ImageListView

urlpatterns = [
    path("v1/", ImageListView.as_view(), name="image-list"),
]
