# Templates

Optional **custom** images when you do not want fat Hub bases.

| Dir | Role |
| --- | --- |
| [kit-core](kit-core/) | Docker `FROM` parent. Not imported into sbx. |
| [kit-shell](kit-shell/) | Minimum empty image you load. Add kits. |
| [kit-cursor](kit-cursor/) | Cursor CLI FROM kit-core. |

```bash
sbx-kit image ls
sbx-kit image load --engine docker kit-shell
sbx-kit image load --engine docker kit-cursor
```

`image load` docker-builds kit-core first. Stock recipes need no pin.
Agent updates are host-side before attach — do not rebake for CLI churn.
Live recipe list: `sbx-kit recipes`.
