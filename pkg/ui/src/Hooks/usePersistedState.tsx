import React, {SetStateAction, useEffect, useState} from "react";

export const usePersistedState =
    <T extends any>({key, defaultValue}: { key: string, defaultValue: T }
    ): [T, React.Dispatch<SetStateAction<T>>] => {
        const [val, setVal] = useState<T>(defaultValue);

        useEffect(() => {
            const sv = localStorage.getItem(key)
            if (sv !== null) setVal(JSON.parse(sv) as T);
            else localStorage.setItem(key, JSON.stringify(defaultValue))
        }, [])

        useEffect(() => {
            localStorage.setItem(key, JSON.stringify(val))
        })

        return [val, setVal as React.Dispatch<SetStateAction<T>>]
    }