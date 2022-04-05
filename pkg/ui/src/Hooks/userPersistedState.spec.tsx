import {act, renderHook} from "@testing-library/react-hooks";
import {usePersistedState} from "./usePersistedState";


describe("usePersistedState", () => {
    const key = "randomKey"
    const defaultValue = 1
    let result: { current: any };
    beforeEach(() => {
        result = renderHook(() => usePersistedState<number>({key, defaultValue})).result
    })
    it("should render the hook correctly", () => {
        expect(result.current).not.toBeUndefined()
    })
    it("should set the default value correctly", () => {
        const [val, _] = result.current
        expect(val).toEqual(defaultValue)
    })
    it("should set the current value correctly", () => {
        const [val, setVal] = result.current
        act(() => {
            setVal(2)
        })
        const {result: resultTwo} = renderHook(() => usePersistedState<number>({key, defaultValue}))
        expect(resultTwo.current[0]).not.toEqual(defaultValue)
        expect(resultTwo.current[0]).toEqual(2)
    })
})