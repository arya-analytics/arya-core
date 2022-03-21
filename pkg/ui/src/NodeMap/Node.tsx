import { motion } from "framer-motion";
import {calculateProgress, HexagonProgressProps, MultiHexagonProgress,} from "./Hexagon";
import React from "react";
import {Typography} from "@mui/material";

export interface NodeMetric {
    key: string;
    value: number;
    max: number;
}

export interface Node {
    id: number;
    metrics: NodeMetric[];
}

export interface NodePinProps {
    node: Node;
    center: [number, number];
}

const PROGRESS_STROKE_WIDTH = 10;
const HEX_EDGE_LENGTH = 60;

const NODE_METRIC_COLORS = ["#49AA19", "#FF008C", "#3C86E3"]

export const NodePin = ({node: {metrics, id}, center = [0, 0]}: NodePinProps) => {
    const progress = metrics.map(({key, value, max}, i): HexagonProgressProps => {
        return {
            strokeWidth: PROGRESS_STROKE_WIDTH,
            progress: calculateProgress(value, max),
            stroke: NODE_METRIC_COLORS[i],
        };
    });
    return (
        <motion.svg x={center[0]} y={center[1]} >
            <motion.text fontSize={48} x={75} y={105} fontFamily={"roboto"} >{id}</motion.text>
            <MultiHexagonProgress progress={progress} edgeLength={HEX_EDGE_LENGTH} center={[90, 90]}/>
        </motion.svg>
    );
};
