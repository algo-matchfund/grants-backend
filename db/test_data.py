#!/usr/bin/python3

import os.path
from ssl import cert_time_to_seconds
import string
import random
import sys


class Project:
    def __init__(
        self,
        name,
        region,
        vat_number,
        address,
        icon,
        description,
        created_by,
        coc_number,
        content,
        app_id,
    ):
        self.name = name
        self.region = region
        self.vat_number = vat_number
        self.address = address
        self.icon = icon
        self.description = description
        self.created_by = created_by
        self.coc_number = coc_number
        self.content = content
        self.app_id = app_id


class MatchingRound:
    def __init__(self, start_date, end_date, match_amount):
        self.start_date = start_date
        self.end_date = end_date
        self.match_amount = match_amount


def loadData(base_dir):
    data_dir = f"{base_dir}/data"
    project_file = f"{data_dir}/projects.csv"

    if not os.path.isdir(data_dir):
        print(f"data dir does not exist")

    # open content file
    content = ""
    with open(f"{base_dir}/test_content.md") as md:
        content = md.read()

    projects = []
    with open(project_file) as f:
        lines = f.read().splitlines()
        for i in range(0, len(lines)):
            r = lines[i]
            name = r.split(",")[0]
            region = r.split(",")[1]
            address = r.split(",")[2]
            icon = r.split(",")[3]
            description = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum gravida dui non est fermentum convallis."
            created_by = "f5b0e5d8-f039-4ea9-aa2b-00f5f10ff77a"
            app_id = i + 1 #TODO: use real app id from deployed apps
            projects.append(
                Project(
                    name,
                    region,
                    f"{i+1:06}",
                    address,
                    icon,
                    description,
                    created_by,
                    i + 1,
                    content,
                    app_id,
                )
            )

    return projects


def truncate_all(f):
    f.write("-- Truncate all tables...\n")
    tables = ["projects", "users", "user_roles", "matching_rounds"]

    for t in tables:
        f.write(f"TRUNCATE {t} RESTART IDENTITY CASCADE;\n")
    f.write("\n")


def insert_projects(f, projects):
    f.write("-- Insert projects...\n")
    for p in projects:
        # project testnet wallet seed: doctor lift possible inspire awake hurry reform unknown illness funny abuse adapt dry strategy leopard neck punch soccer differ clown surface blossom crash abandon rigid
        f.write(
            f"INSERT INTO PROJECTS (name, algorand_wallet, icon, description, created_by, content, app_id) "
            f"VALUES ('{p.name}', 'VO6X5ACR5HOPS3J3CGBVVKLWYDD3UWHPXWQCM22Q4N7GOLXCXNPLAKAA7A', '{p.icon}', '{p.description}', '{p.created_by}', '{p.content}', {p.app_id});\n"
        )

    f.write("\n")


def insert_matching_round(f):
    f.write("-- Insert matching round...\n")
    f.write(
        f"INSERT INTO MATCHING_ROUNDS (start_date, end_date, match_amount) "
        f"VALUES ('2022-03-01 00:00:00', '2023-05-31 23:59:59', '100000000000');\n"
    )
    f.write("\n")


def insert_users(f):
    # For now this just inserts users that are seeded in keycloak in our deployment
    user_insert = (
        "INSERT INTO users (id, github_authorized, is_company) values ('{}', {}, {});\n"
    )
    # role_insert = "INSERT INTO user_roles (user_id, role) values ('{}', '{}');\n"

    # ## Insert our demo TA Admin user
    # f.write(user_insert.format("e8f78c06-fe2c-471f-9451-b9111041aada", "admin@ta.org"))
    # f.write(role_insert.format("e8f78c06-fe2c-471f-9451-b9111041aada", "admin"))

    # ## Insert our demo TA Admin user
    # f.write(user_insert.format("ac8aa240-34d5-4f8c-82b6-30698edd4a42", "user@ta.org"))
    # f.write(role_insert.format("ac8aa240-34d5-4f8c-82b6-30698edd4a42", "user"))

    ## Insert our demo Company User
    f.write(user_insert.format("f5b0e5d8-f039-4ea9-aa2b-00f5f10ff77a", "FALSE", "TRUE"))
    f.write(
        user_insert.format("e8f78c06-fe2c-471f-9451-b9111041aada", "FALSE", "fALSE")
    )
    f.write(
        user_insert.format("ac8aa240-34d5-4f8c-82b6-30698edd4a42", "FALSE", "fALSE")
    )
    f.write(
        user_insert.format("d29e2a39-cff6-4447-a33b-12b82ad5cffd", "FALSE", "fALSE")
    )

    # f.write(role_insert.format("660ca053-3c41-414b-8fa8-0d47760fb208", "user"))


def insert_funding_campaigns(f, projects):
    f.write("-- Insert funding campaigns...\n")
    for i in range(0, len(projects)):
        project_id = i + 1
        start_time = '2022-03-01 00:00:00'
        f.write(
            f"INSERT INTO FUNDING_CAMPAIGNS (project_id, start_time) values ('{project_id}', '{start_time}');\n"
        )

    f.write("\n")


def insert_contributions(f, projects):
    f.write("-- Insert funding campaigns...\n")
    for i in range(0, len(projects)):
        project_id = 1 + i
        user_id = "e8f78c06-fe2c-471f-9451-b9111041aada"
        amount = random.randint(10000, 10000000)
        f.write(
            "INSERT INTO CONTRIBUTIONS (project_id, matching_round_id, user_id, amount) "
            f"values ({project_id}, 1, '{user_id}', {amount});\n"
        )

    for i in range(0, len(projects)):
        project_id = 1 + i
        user_id = "ac8aa240-34d5-4f8c-82b6-30698edd4a42"
        amount = random.randint(10000, 10000000)
        f.write(
            "INSERT INTO CONTRIBUTIONS (project_id, matching_round_id, user_id, amount) "
            f"values ({project_id}, 1, '{user_id}', {amount});\n"
        )

    f.write("\n")


def generate(projects, base_dir):
    f = open(f"{base_dir}/test_data.sql", "w")

    truncate_all(f)
    insert_users(f)
    insert_projects(f, projects)
    insert_matching_round(f)
    insert_funding_campaigns(f, projects)
    insert_contributions(f, projects)

    f.close()


def main():
    print("usage: ./test_data.py [<data dir>]")

    base_dir = "."
    if len(sys.argv) == 2:
        base_dir = sys.argv[1]

    print(f"Generating data...")
    projects = loadData(base_dir)
    generate(projects, base_dir)
    print(f"Finished generating data...")


if __name__ == "__main__":
    main()
