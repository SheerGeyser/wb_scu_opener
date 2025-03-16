#Requires AutoHotkey v2.0

^+f:: {
    initialBuffer := A_Clipboard
    Send("^c")
    Sleep(100)
    copiedText := A_Clipboard

    if !copiedText {
        return
    }

    matches := []
    pos := 1

    while pos := RegExMatch(copiedText, "\b\d{4,}\b", &match, pos) {
        matches.Push(match[0])
        pos += StrLen(match[0])
    }


    for id in matches {
        Run(Format("https://www.wildberries.ru/catalog/{}/detail.aspx", id))
        Sleep(200)
    }


    A_Clipboard := initialBuffer
    Sleep(100)
}
