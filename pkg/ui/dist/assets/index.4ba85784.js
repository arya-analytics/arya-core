var y = Object.defineProperty;
var u = Object.getOwnPropertySymbols;
var w = Object.prototype.hasOwnProperty, v = Object.prototype.propertyIsEnumerable;
var p = (e, t, r) => t in e ? y(e, t, {enumerable: !0, configurable: !0, writable: !0, value: r}) : e[t] = r,
    l = (e, t) => {
        for (var r in t || (t = {})) w.call(t, r) && p(e, r, t[r]);
        if (u) for (var r of u(t)) v.call(t, r) && p(e, r, t[r]);
        return e
    };
import {
    a as k,
    b as c,
    B as h,
    c as T,
    d,
    e as C,
    f as P,
    g,
    h as I,
    i as M,
    j as o,
    m as E,
    r as a,
    R as B,
    S as L,
    T as S,
    u as b
} from "./vendor.de4270d2.js";

const N = function () {
    const t = document.createElement("link").relList;
    if (t && t.supports && t.supports("modulepreload")) return;
    for (const n of document.querySelectorAll('link[rel="modulepreload"]')) i(n);
    new MutationObserver(n => {
        for (const s of n) if (s.type === "childList") for (const m of s.addedNodes) m.tagName === "LINK" && m.rel === "modulepreload" && i(m)
    }).observe(document, {childList: !0, subtree: !0});

    function r(n) {
        const s = {};
        return n.integrity && (s.integrity = n.integrity), n.referrerpolicy && (s.referrerPolicy = n.referrerpolicy), n.crossorigin === "use-credentials" ? s.credentials = "include" : n.crossorigin === "anonymous" ? s.credentials = "omit" : s.credentials = "same-origin", s
    }

    function i(n) {
        if (n.ep) return;
        n.ep = !0;
        const s = r(n);
        fetch(n.href, s)
    }
};
N();
const O = e => a.exports.createElement("svg", l({
        xmlns: "http://www.w3.org/2000/svg",
        viewBox: "0 0 200 189.34",
        width: "1em",
        height: "1em"
    }, e), a.exports.createElement("defs", null, a.exports.createElement("style", null, ".cls-1{fill:#fff;}")), a.exports.createElement("g", {
        id: "Layer_2",
        "data-name": "Layer 2"
    }, a.exports.createElement("g", {
        id: "Layer_1-2",
        "data-name": "Layer 1"
    }, a.exports.createElement("path", {
        className: "cls-1",
        d: "M52.61,168.82A14.81,14.81,0,0,1,39.78,146.6L94.22,52.33c5.7-9.88,20-9.88,25.9.42l51.77,89.67a6.88,6.88,0,0,0,2.48,2.49l15.42,8.9a6.78,6.78,0,0,0,9.26-9.27L119.87,7.41a14.8,14.8,0,0,0-25.65,0L2,167.12a14.81,14.81,0,0,0,12.83,22.22H170.39a6.79,6.79,0,0,0,3.39-12.66l-12.05-7a6.83,6.83,0,0,0-3.39-.91Z"
    })))), z = e => a.exports.createElement("svg", l({
        xmlns: "http://www.w3.org/2000/svg",
        xmlnsXlink: "http://www.w3.org/1999/xlink",
        viewBox: "0 0 200 189.34",
        width: "1em",
        height: "1em"
    }, e), a.exports.createElement("defs", null, a.exports.createElement("style", null, ".cls-1{opacity:0.85;fill:url(#linear-gradient);}"), a.exports.createElement("linearGradient", {
        id: "linear-gradient",
        x1: -2.63,
        y1: 56.85,
        x2: 167.19,
        y2: 157.88,
        gradientUnits: "userSpaceOnUse"
    }, a.exports.createElement("stop", {offset: 0, stopColor: "#10007f"}), a.exports.createElement("stop", {
        offset: 1,
        stopColor: "#0084e9"
    }))), a.exports.createElement("g", {
        id: "Layer_2",
        "data-name": "Layer 2"
    }, a.exports.createElement("g", {
        id: "Layer_1-2",
        "data-name": "Layer 1"
    }, a.exports.createElement("path", {
        className: "cls-1",
        d: "M52.61,168.82A14.81,14.81,0,0,1,39.78,146.6L94.22,52.33c5.7-9.88,20-9.88,25.9.42l51.77,89.67a6.88,6.88,0,0,0,2.48,2.49l15.42,8.9a6.78,6.78,0,0,0,9.26-9.27L119.87,7.41a14.8,14.8,0,0,0-25.65,0L2,167.12a14.81,14.81,0,0,0,12.83,22.22H170.39a6.79,6.79,0,0,0,3.39-12.66l-12.05-7a6.83,6.83,0,0,0-3.39-.91Z"
    })))), _ = e => {
        const {palette: {mode: t}} = b(), r = t == "light" ? e.light : e.dark;
        return o(L, l({component: r, inheritViewBox: !0}, e))
    }, F = e => o(_, l({light: z, dark: O}, e)), R = ({key: e, defaultValue: t}) => {
        const [r, i] = a.exports.useState(t);
        return a.exports.useEffect(() => {
            const n = localStorage.getItem(e);
            n !== null ? i(JSON.parse(n)) : localStorage.setItem(e, JSON.stringify(t))
        }, []), a.exports.useEffect(() => {
            localStorage.setItem(e, JSON.stringify(r))
        }), [r, i]
    }, A = "dark", $ = "aryaTheme", f = a.exports.createContext({
        palette: "light", theme: {}, setPalette: () => {
        }
    }), j = () => a.exports.useContext(f), H = ({children: e}) => {
        const [t, r] = R({key: $, defaultValue: A}), i = D(t);
        return a.exports.useEffect(() => {
            var n, s;
            console.log(i), document.body.style.backgroundColor = (s = (n = i.palette) == null ? void 0 : n.background) == null ? void 0 : s.paper
        }, [t]), o(S, {theme: i, children: o(f.Provider, {value: {palette: t, theme: i, setPalette: r}, children: e})})
    }, D = e => K[e], J = {
        palette: {primary: {main: "#3774D0"}},
        shape: {borderRadius: 2},
        components: {
            MuiTabs: {
                defaultProps: {
                    sx: {
                        borderBottom: 1,
                        borderColor: "divider",
                        minHeight: 0,
                        "& .MuiTabs-indicator": {backgroundColor: ""},
                        "& .MuiButtonBase-root": {height: 36, minHeight: 0, textTransform: "none", color: "text.primary"}
                    }
                }
            }, MuiTypography: {defaultProps: {color: "text.primary"}}
        }
    }, U = {palette: {mode: "light", secondary: {main: "#212121"}, text: {primary: "#212121"}}}, G = {
        palette: {
            mode: "dark",
            background: {default: "#1F1F1F"},
            secondary: {main: "#e0e0e0"},
            text: {primary: "#e0e0e0"}
        }
    }, x = e => T(E(J, e)), K = {light: x(U), dark: x(G)}, V = e => {
        const {palette: t, setPalette: r} = j();
        return o(k, l({onChange: () => r(t == "light" ? "dark" : "light"), checked: t == "light"}, e))
    }, Z = e => ({hidden: {width: 0}, visible: t => ({width: `${e}%`})}),
    q = ({progress: e, fill: t, name: r, width: i}) => c(d, {
        direction: "row",
        spacing: 2,
        sx: {display: "flex", alignItems: "center"},
        children: [o(C, {variant: "subtitle2", children: r}), o(h, {
            sx: {
                width: i,
                height: "10px",
                backgroundColor: "action.disabledBackground",
                margin: "5px 0"
            },
            children: o(P.div, {
                initial: "hidden",
                animate: "visible",
                variants: Z(e),
                style: {width: `${e}%`, height: "10px", backgroundColor: t}
            })
        })]
    }), X = ({}) => c(h, {
        sx: {
            width: "100vw",
            height: "100vh",
            display: "flex",
            justifyContent: "center",
            alignItems: "center"
        }, children: [o(W, {}), o(Y, {}), o(Q, {})]
    }), Q = () => c(d, {
        direction: "column",
        spacing: 6,
        sx: {display: "flex", alignItems: "center", marginTop: -7},
        children: [o(F, {fontSize: "large", sx: {fontSize: 150}}), c(d, {
            direction: "column",
            spacing: 4,
            children: [o(g, {
                label: "Username",
                sx: {width: 460},
                size: "small",
                variant: "standard"
            }), o(g, {label: "Password", sx: {width: 460}, size: "small", variant: "standard"})]
        }), o(I, {variant: "contained", size: "medium", children: "Login"})]
    }), W = () => o(V, {size: "small", sx: {position: "absolute", bottom: "0", left: "0", zIndex: "1", m: "13px"}}),
    Y = () => o(d, {
        direction: "row",
        sx: {position: "absolute", bottom: 0, right: 0, zIndex: 1, m: "13px"},
        children: o(q, {name: "Live Nodes", progress: 88, fill: "green", width: 400})
    });

function ee() {
    return o(h, {sx: {flexGrow: 1, height: "100vh"}, children: o(X, {})})
}

B.render(o(M.StrictMode, {children: o(H, {children: o(ee, {})})}), document.getElementById("root"));
