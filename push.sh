#!/bin/sh

repository="https://Kajekk:$github_token@github.com/Kajekk/useful-repositories.git"
username="Kajekk"
# Git Email
email="tangquocvinh1000@gmail.com"
# Branch
branch="main"
# Go file
go_script="main.go"

current_time=$(date +"%Y-%m-%dT%H:%M:%S")

cd ~/workspace/useful-repositories

git config user.name $username
git config user.email $email
git remote set-url origin $repository
git fetch origin $branch
git reset origin/$branch
git checkout $branch

go run $go_script

rm README.md
mv README_TEMP.md README.md

git add -A .
git commit -m "Auto update at $current_time"
git push -q origin $branch

read