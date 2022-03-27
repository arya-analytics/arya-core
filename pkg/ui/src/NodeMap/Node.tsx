import { motion } from "framer-motion";
import {
  calculateProgress,
  Hexagon,
  HexagonProgressProps,
  MultiHexagonProgress,
} from "./Hexagon";
import React from "react";
import { useThemeContext } from "../Theme/Theme";
import { CanvasItemProps } from "../Canvas/Canvas";

export interface NodeMetric {
  key: string;
  value: number;
  max: number;
}

export interface Node {
  id: number;
  metrics: NodeMetric[];
}

export interface NodePinProps extends Partial<CanvasItemProps> {
  node: Node;
  onClick?: (node: Node) => void;
}

const PROGRESS_STROKE_WIDTH = 10;
const HEX_EDGE_LENGTH = 60;

export const NODE_METRIC_COLORS = [
  "#49AA19",
  "#FF008C",
  "#3C86E3",
  "#ee970a",
  "#a118fa",
];

export const NodePin = ({ node, onClick, ...props }: NodePinProps) => {
  const { metrics, id } = node;
  const { theme } = useThemeContext();
  const progress = metrics
    .slice(0, 3)
    .map(({ key, value, max }, i): HexagonProgressProps => {
      return {
        name: key,
        strokeWidth: PROGRESS_STROKE_WIDTH,
        progress: calculateProgress(value, max),
        stroke: NODE_METRIC_COLORS[i],
      };
    });
  return (
    <motion.svg
      {...props}
      cursor={"pointer"}
      opacity={1}
      whileTap={{ opacity: 0.7 }}
      whileHover={{ opacity: 0.85 }}
      onMouseDown={(e) => {
        if (onClick) {
          onClick(node);
        }
      }}
    >
      <Hexagon edgeLength={HEX_EDGE_LENGTH} center={[90, 90]} opacity={0} />
      <motion.text
        fontSize={48}
        x={75}
        y={105}
        fontFamily={"roboto"}
        fill={theme.palette?.text?.primary}
        transform={props.transform}
      >
        {id}
      </motion.text>
      <MultiHexagonProgress
        progress={progress}
        edgeLength={HEX_EDGE_LENGTH}
        center={[90, 90]}
        transform={props.transform}
      />
    </motion.svg>
  );
};
