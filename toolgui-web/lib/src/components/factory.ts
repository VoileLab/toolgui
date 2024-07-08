import { TTextbox } from "./tcinput/textbox"
import { TCheckbox } from "./tcinput/checkbox"
import { TButton } from "./tcinput/button"
import { TSelect } from "./tcinput/select"
import { TTextarea } from "./tcinput/textarea"
import { TFileupload } from "./tcinput/fileupload"
import { TRadio } from "./tcinput/radio"
import { TDatepicker } from "./tcinput/datepicker"

import { TContainer } from "./tclayout/container"
import { TBox } from "./tclayout/box"
import { TColumn } from "./tclayout/column"

import { TTitle } from "./tccontent/title"
import { TImage } from "./tccontent/image"
import { TSubtitle } from "./tccontent/subtitle"
import { TText } from "./tccontent/text"
import { TDivider } from "./tccontent/divider"
import { TMarkdown } from "./tccontent/markdown"
import { TCode } from "./tccontent/code"

import { TJson } from "./tcdata/json"
import { TTable } from "./tcdata/table"

import { TProgressar } from "./tcmisc/progress_bar"
import { TMessage } from "./tcmisc/message"

import { Props } from "./component_interface"
import { TLink } from "./tccontent/link"

const creatorMap: { [id: string]: ((props: Props) => JSX.Element) } = {
  textbox_component: TTextbox,
  checkbox_component: TCheckbox,
  button_component: TButton,
  select_component: TSelect,
  textarea_component: TTextarea,
  fileupload_component: TFileupload,
  radio_component: TRadio,
  datepicker_component: TDatepicker,

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
  link_component: TLink,

  json_component: TJson,
  table_component: TTable,

  progress_bar_component: TProgressar,
  message_component: TMessage,
}


export function TComponent({ node, update, nodes }: Props) {
  if (!(node.props.name in creatorMap)) {
    throw new Error(`unsupported component type: ${node.props.name}`);
  }

  return creatorMap[node.props.name]({ node, update, nodes })
}
