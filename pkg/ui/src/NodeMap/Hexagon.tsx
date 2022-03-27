import { motion, SVGMotionProps } from "framer-motion";
import React from "react";
import { useThemeContext } from "../Theme/Theme";

// |||| BASE ||||

export interface HexagonProps extends SVGMotionProps<SVGPolygonElement> {
  center?: [number, number];
  edgeLength?: number;
}

export const Hexagon = ({
  edgeLength = 10,
  center = [0, 0],
  ...props
}: HexagonProps) => {
  return (
    <motion.polygon points={hexPoints({ edgeLength, center })} {...props} />
  );
};

// |||| PROGRESS ||||

export const calculateProgress = (value: number, max: number) => value / max;

export interface HexagonProgressProps extends HexagonProps {
  name?: string;
  strokeWidth?: number;
  progress: number;
}

export const HexagonProgress = ({
  progress,
  ...props
}: HexagonProgressProps) => {
  return (
    <Hexagon
      initial="hidden"
      animate="visible"
      fill="none"
      stroke={"blue"}
      variants={drawProgress(progress)}
      {...props}
    />
  );
};

export interface MultiHexagonProgress
  extends Omit<HexagonProgressProps, "progress"> {
  progress: HexagonProgressProps[];
}

export const MultiHexagonProgress = ({
  progress: progressValues,
  drag,
  edgeLength = 10,
  progress,
  strokeWidth = 10,
  center = [0, 0],
  ...props
}: MultiHexagonProgress) => {
  const { theme } = useThemeContext();
  return (
    <>
      {progressValues.map(({ progress, edgeLength: _, ...hexProps }, i) => {
        return (
          <HexagonProgress
            id={`${i}`}
            progress={progress}
            edgeLength={edgeLength + i * strokeWidth}
            strokeWidth={strokeWidth}
            center={center}
            {...hexProps}
            {...props}
          />
        );
      })}
    </>
  );
};

// |||| PROGRESS ANIMATION ||||

const PROGRESS_ANIMATION_DURATION = 0.5;
const PROGRESS_ANIMATION_TYPE = "sprint";
const PROGRESS_ANIMATION_BOUNCE = 0;

const drawProgress = (progress: number) => {
  return {
    hidden: { pathLength: 0 },
    visible: (i: number) => {
      return {
        pathLength: progress,
        transition: {
          pathLength: {
            type: PROGRESS_ANIMATION_TYPE,
            duration: PROGRESS_ANIMATION_DURATION,
            bounce: PROGRESS_ANIMATION_BOUNCE,
          },
        },
      };
    },
  };
};

// |||| POINT GENERATION ||||

const HEX_COS = Math.abs(Math.cos((2 * Math.PI) / 3));
const HEX_SIN = Math.abs(Math.sin((2 * Math.PI) / 3));

const hexPoints = ({
  edgeLength,
  center = [0, 0],
}: {
  edgeLength: number;
  center: [number, number];
}): string => {
  const cw = HEX_SIN * edgeLength;
  const sw = HEX_COS * edgeLength;
  var ptStr = "";
  const vals = [
    [0, edgeLength / 2 + sw],
    [-cw, edgeLength / 2],
    [-cw, -edgeLength / 2],
    [0, -(edgeLength / 2 + sw)],
    [cw, -edgeLength / 2],
    [cw, edgeLength / 2],
  ];
  vals
    .map(([x, y]) => [x + center[0], y + center[1]])
    .forEach(([x, y]) => (ptStr += `${x} ${y},`));

  return ptStr;
};
