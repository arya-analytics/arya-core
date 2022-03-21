import React from "react";
import {Node, NodePin} from "./Node";
import {Button, ButtonGroup, Grid, Stack, Typography} from "@mui/material";
import {FitScreenOutlined, NoiseAwareOutlined} from "@mui/icons-material";
import {motion} from "framer-motion";
import {useThemeContext} from "../Theme/Theme";


const DUMMY_NODES: Node[] = [
    {
        id: 1,
        metrics: [
            {
                key: "Memory",
                value: 12,
                max: 32,
            },
            {
                key: "CPU",
                value: 4.32,
                max: 8.0,
            },
            {
                key: "Storage",
                value: 300,
                max: 900,
            }
        ],
    },
    {
        id: 2,
        metrics: [
            {
                key: "Memory",
                value: 12,
                max: 16,
            },
            {
                key: "CPU",
                value: 7.32,
                max: 8.0,
            },
            {
                key: "Storage",
                value: 450,
                max: 900,
            }
        ],
    },
    {
        id: 3,
        metrics: [
            {
                key: "Memory",
                value: 96,
                max: 128,
            },
            {
                key: "CPU",
                value: 4.32,
                max: 12.0,
            },
            {
                key: "Storage",
                value: 722,
                max: 900,
            }
        ],
    },
];

export const NodeMapPage = () => {
    const nodes = DUMMY_NODES;
    return (
        <Grid
            xs={8}
            borderRight="1px solid"
            borderColor="divider"
            direction="column"
            container
        >
            <Grid
                height={30}
                borderBottom="1px solid"
                borderColor="divider"
                paddingLeft={1}
                display="flex"
                alignItems="center"
                justifyContent="space-between"
            >
                <Stack direction="row" alignItems="center" spacing={1}>
                    <NoiseAwareOutlined
                        sx={{color: "text.primary"}}
                        fontSize="small"
                    />
                    <Typography variant="subtitle1">Node Map</Typography>
                </Stack>
            </Grid>
            <NodeMap nodes={nodes}/>
        </Grid>
    )
}

const ONE_NODE_CENTER = [
    [0.3, 0.35]
]

const TWO_NODE_CENTERS = [
    [0.2, 0.35],
    [0.4, 0.35]
]

const THREE_NODE_CENTERS = [
    [0.3, 0.1],
    [0.2, 0.4],
    [0.4, 0.4],
]

const NodePatterns: { [key: string]: any } = {
    "1": ONE_NODE_CENTER,
    "2": TWO_NODE_CENTERS,
    "3": THREE_NODE_CENTERS,
}

const scaleNodePatterns = (width: number, height: number, pattern: [number, number]): [number, number] => {
    return [pattern[0] * width, pattern[1] * height];
}

const NodeMap = ({nodes}: { nodes: Node[] }) => {
    const nodePattern = NodePatterns[nodes.length.toString()];
    const [theme, _] = useThemeContext()
    return (
        <Grid
            xs
            sx={{
                backgroundColor: theme == "light" ? "grey.300" : "grey.900",
                overflow: "hidden",
                position: "relative",
                cursor: "pointer"
            }}
            item
        >
            <ButtonGroup variant="outlined" size="small" sx={{position: "absolute", right: 15, top: 15}}
                         color="secondary">
                <Button><FitScreenOutlined/></Button>
                <Button>Two</Button>
                <Button>Three</Button>
            </ButtonGroup>
            <motion.svg width="100%" height="100%"
                        onPan={(v, p) => console.log({v, p})}
            >
                {nodes.map((n, i) => {
                    const center = scaleNodePatterns(window.innerWidth, window.innerHeight, nodePattern[i])
                    return (<NodePin node={n} center={center}/>)
                })}
            </motion.svg>
        </Grid>
    );
};
