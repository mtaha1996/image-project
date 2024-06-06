from django.db import migrations


def load_image_data(apps, schema_editor):
    Image = apps.get_model("image_bank", "Image")
    data = [
        {
            "query": "cute kittens",
            "link": "https://images.unsplash.com/photo-1614035030394-b6e5b01e0737",
        },
        {
            "query": "cute kittens",
            "link": "https://images.unsplash.com/photo-1556582305-528bffcf7af0",
        },
    ]

    for item in data:
        Image.objects.create(query=item["query"], link=item["link"])

    Access = apps.get_model("image_bank", "Access")
    data = [
        {
            "engin_id": "engin",
            "api_key": "key",
        },
    ]

    for item in data:
        Access.objects.create(engin_id=item["engin_id"], api_key=item["api_key"])


class Migration(migrations.Migration):

    dependencies = [
        ("image_bank", "0001_initial"),
    ]

    operations = [
        migrations.RunPython(load_image_data),
    ]
