import { TTextbox } from "./textbox"
import { TCheckbox } from "./checkbox"
import { TButton } from "./button"
import { TSelect } from "./select"

import { TContainer } from "./layouts"
import { TBox } from "./layouts"
import { TColumn } from "./layouts"

import { TTitle } from "./contents"
import { TImage } from "./contents"
import { TSubtitle } from "./contents"
import { TText } from "./contents"
import { TDivider } from "./contents"
import { TMarkdown } from "./contents"
import { TCode } from "./contents"

import { TMessage } from "./message"
import { TJson } from "./json"

import { TProgressar } from "./progress_bar"

const creatorMap = {
  textbox_component: TTextbox,
  checkbox_component: TCheckbox,
  button_component: TButton,
  select_component: TSelect,

  container_component: TContainer,
  box_component: TBox,
  column_component: TColumn,

  title_component: TTitle,
  subtitle_component: TSubtitle,
  image_component: TImage,
  text_component: TText,
  divider_component: TDivider,
  markdown_component: TMarkdown,
  code_component: TCode,

  message_component: TMessage,
  json_component: TJson,

  progress_bar_component: TProgressar,
}

export function TComponent({ node, update, nodes }) {
  return creatorMap[node.props.name]({ node, update, nodes })
}
