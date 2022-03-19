import "./App.tsx"
import App from "./App";
import {render} from "@testing-library/react"

describe("test", () => {
    it("should be true", () => {
        render(<App />)
        expect(true).toBe(true)
    })
})