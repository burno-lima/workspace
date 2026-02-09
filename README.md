# 🚀 Workspace - Project Workspace Manager TUI

Um gerenciador de workspace de projetos no estilo lazygit feito em Go com interface TUI (Terminal User Interface).

## ✨ Funcionalidades

- 📂 Lista todos os diretórios do path especificado
- 🎨 Interface colorida e moderna estilo lazygit
- 📊 Layout em colunas (como `ls`) para melhor visualização
- ⌨️  Navegação intuitiva com setas ou vim keybindings (hjkl)
- ✅ Abre o projeto selecionado no editor configurado com Enter
- 📍 Mostra informações do projeto selecionado
- 🎯 Interface responsiva que se adapta ao tamanho do terminal
- 🔍 Busca rápida de projetos com `/`
- ⚙️  Configuração de editor padrão e específico por projeto
- 🎨 Ícones de tipo de projeto (Go, Python, JavaScript, etc.) usando glifos Unicode

## 📦 Instalação

```bash
go build -o workspace ./cmd/workspace
```

## 🎨 Fontes e Glifos

Este projeto usa **glifos Unicode** (Nerd Fonts) para exibir ícones no terminal. Para melhor experiência visual, instale uma fonte Nerd Font:

### Fontes Recomendadas:
- [FiraCode Nerd Font](https://www.nerdfonts.com/font-downloads) - Recomendada
- [JetBrains Mono Nerd Font](https://www.nerdfonts.com/font-downloads)
- [Hack Nerd Font](https://www.nerdfonts.com/font-downloads)
- [Meslo Nerd Font](https://www.nerdfonts.com/font-downloads)

### Como instalar (Linux):
```bash
# Baixar fonte
wget https://github.com/ryanoasis/nerd-fonts/releases/download/v3.1.1/FiraCode.zip

# Extrair
unzip FiraCode.zip -d ~/.local/share/fonts/

# Atualizar cache de fontes
fc-cache -fv
```

### Configurar no terminal:
- **Alacritty**: Edite `~/.config/alacritty/alacritty.yml` e adicione:
  ```yaml
  font:
    normal:
      family: "FiraCode Nerd Font"
  ```
- **Kitty**: Edite `~/.config/kitty/kitty.conf`:
  ```
  font_family FiraCode Nerd Font
  ```
- **GNOME Terminal / Tilix**: Nas preferências, selecione a fonte Nerd Font instalada

### Ícones disponíveis:
-  - Pasta/diretório
-  - Busca
-  - Erro
-  - Configuração
-  - Target/Alvo
-  - Editor de texto
-  - IDE (JetBrains)
- ▶ - Seleção/Arrow
-  - Projeto genérico
-  - Go
-  - Python
-  - JavaScript/TypeScript
-  - Java
-  - PHP
-  - Ruby
-  - Rust
-  - C/C++

**Nota**: Se você não tiver uma Nerd Font instalada, os ícones aparecerão como quadrados vazios ou caracteres estranhos. A funcionalidade do programa não será afetada, apenas a estética.

## 🚀 Uso

```bash
./workspace <path>
```

Exemplo:
```bash
./workspace /home/brunolima/workspace/napp
```

## ⌨️ Controles

### Navegação Normal:
- `↑` ou `k`: Mover para cima
- `↓` ou `j`: Mover para baixo
- `←` ou `h`: Mover para a esquerda (coluna anterior)
- `→` ou `l`: Mover para a direita (próxima coluna)
- `Enter`: Abrir projeto selecionado
- `/`: Ativar modo de busca
- `c`: Abrir menu de configuração
- `q` ou `Ctrl+C`: Sair

### Modo de Busca:
- Digite para buscar projetos
- `Backspace`: Apagar caractere
- `Enter`: Abrir projeto filtrado
- `Esc`: Cancelar busca
- `↑↓←→` ou `Ctrl+n/p`: Navegar nos resultados

### Menu de Configuração:
- `↑↓` ou `k/j`: Navegar opções
- `Enter`: Selecionar opção
- `Esc` ou `q`: Fechar menu

## 📋 Requisitos

- Go 1.21 ou superior
- Pelo menos um editor/IDE instalado (VSCode, Cursor, Neovim, etc.)
- (Opcional) Nerd Font para ícones bonitos

## 📚 Dependências

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Framework TUI
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Estilização de terminal
