# loadenv.zsh
# Usage: source loadenv.zsh [path/to/.env]

# Проверяем, что скрипт ИМЕННО sourced (а не исполняется как процесс)
if [[ "${ZSH_EVAL_CONTEXT}" != *:file* ]]; then
  echo "➡️  Use: source loadenv.zsh [path/to/.env]"
  return 1 2>/dev/null || exit 1
fi

# Опции экспорта
setopt allexport

# Определяем путь к скрипту и к .env (по умолчанию рядом со скриптом)
local this="${(%):-%N}"
local script_dir="${this:A:h}"
local envfile="${1:-$script_dir/.env}"

if [[ ! -f "$envfile" ]]; then
  echo "❌ .env not found at: $envfile"
  unsetopt allexport
  return 1
fi

# Загружаем файл
source "$envfile"

# Выключаем авто-export
unsetopt allexport

echo "✅ Loaded: $envfile"