import boto3

client = boto3.client("lambda")

account_number = "131791614471"
ecr_name = "telegram_bot"
image_tag = "latest"
function_name = "telegram_bot"
image_uri = f"{account_number}.dkr.ecr.ap-southeast-1.amazonaws.com/{ecr_name}:{image_tag}"

if __name__ == "__main__":
    response = client.update_function_code(
        FunctionName=function_name,
        ImageUri=image_uri,
    )
    print(response)

