import os
import utils


def main():
    src_file = os.path.join(os.getcwd(), "deploy/postgres/postgres-configmap.yaml")
    dst_file = os.path.join(os.getcwd(), "build/postgres-configmap.yaml")
    with open(src_file, "r") as src:
        with open(dst_file, "w+") as dst:
            data = src.read()
            print("Deploying {}".format(dst_file))
            dst.write(data)

    utils.apply(dst_file)

    src_file = os.path.join(os.getcwd(), "deploy/postgres/postgres-deployment.yaml")
    dst_file = os.path.join(os.getcwd(), "build/postgres-deployment.yaml")
    with open(src_file, "r") as src:
        with open(dst_file, "w+") as dst:
            data = src.read()
            print("Deploying {}".format(dst_file))
            dst.write(data)
    utils.apply(dst_file)

    src_file = os.path.join(os.getcwd(), "deploy/postgres/postgres-storage.yaml")
    dst_file = os.path.join(os.getcwd(), "build/postgres-storage.yaml")
    with open(src_file, "r") as src:
        with open(dst_file, "w+") as dst:
            data = src.read()
            try:
                size = utils.check_output(
                    "kubectl -n assisted-installer get persistentvolumeclaims postgres-pv-claim " +
                    "-o=jsonpath='{.status.capacity.storage}'")
                print("Using existing disk size", size)
            except:
                size = "10Gi"
                print("Using default size", size)
            data = data.replace("REPLACE_STORAGE", size)
            print("Deploying {}".format(dst_file))
            dst.write(data)

    utils.apply(dst_file)


if __name__ == "__main__":
    main()
