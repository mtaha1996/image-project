from django.shortcuts import render
from .permissions import IsValidApiKeyAndEnginId
from rest_framework import generics, status
from .models import Image
from .serializers import ImageSerializer
from rest_framework.exceptions import ValidationError
from rest_framework.response import Response


class ImageListView(generics.ListAPIView):
    permission_classes = [IsValidApiKeyAndEnginId]
    serializer_class = ImageSerializer

    def get_queryset(self):
        query = self.request.query_params.get("q")
        return Image.objects.filter(query=query)

    def list(self, request, *args, **kwargs):
        num = self.request.query_params.get("num")
        try:
            num = int(num)
            if num < 1:
                raise ValidationError("The 'num' parameter must be a positive integer.")
        except (TypeError, ValueError):
            raise ValidationError("The 'num' parameter must be a valid integer.")

        queryset = self.get_queryset()[:num]
        serializer = self.get_serializer(queryset, many=True)
        return Response({"items": serializer.data})
