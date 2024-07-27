#!/bin/bash

# Chemin du fichier d'environnement
env_file="$(dirname "$0")/../integration.env"

# Vérifier si le fichier d'environnement existe
if [ ! -f "$env_file" ]; then
  echo "Le fichier $env_file est manquant"
  exit 1
fi

# Charger les variables d'environnement
echo "Chargement des variables d'environnement depuis $env_file"
source "$env_file"

# Fonction pour exécuter un script et vérifier son code de sortie
run_script() {
  local script="$1"
  echo "Running $script..."
  bash "$script"
  local exit_code=$?
  if [ $exit_code -ne 0 ]; then
    echo "Erreur lors de l'exécution de $script avec le code de sortie $exit_code"
    exit $exit_code
  fi
  echo
}

# Exécuter les scripts dans l'ordre approprié
run_script ./tests/integration/create-api-user.sh
run_script ./tests/integration/create-api-key.sh

# Exécuter les autres scripts d'intégration
for script in ./tests/integration/*.sh; do
    if [[ "$script" != *"create-api-user.sh"* && "$script" != *"create-api-key.sh"* ]]; then
        run_script "$script"
    fi
done

echo "Tous les tests d'intégration ont été exécutés avec succès."
