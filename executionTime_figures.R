library(ggplot2)
library(stringr)

# Your data (replace with your actual data) - Mixed units!
#data_strings <- c("11.063µs", "136ns", "15.66µs", "11.455µs", "10.86µs", "22.755µs", "14.858µs", "2s", "500ms", "10ns", "0.001s", "750µs")
setwd("~/worktable/errorutils")
file_path <- "/home/pollo/worktable/errorutils/timing.log"  #Replace with your file path
# use when newline list: data_strings <- readLines(file_path)
df <- read.csv(file_path,
               colClasses = c("character", "factor"),
               header = F)
data_strings <- df[["V1"]]
readFast_into_ns <- function(time_str) {
  # Define a lookup table for the units and their corresponding multipliers
  unit_multipliers <- c(
    "ns" = 1,
    "µs" = 1e3,
    "ms" = 1e6,
    "s"  = 1e9,
    "µ"  = 1e3
  )
  
  # Extract the unit from the end of the string
  unit <- str_extract(time_str, "ns$|µs$|ms$|s$|µ$")
  
  if (!is.na(unit)) {
    # Remove the unit from the string and convert the remaining part to numeric
    value <- as.numeric(str_replace(time_str, unit, ""))
    # Multiply the value by the corresponding multiplier
    value * unit_multipliers[unit]
  } else {
    NA  # Handle unknown units
  }
}
read_into_ns <- function(time_str) {
  if (str_detect(time_str, "ns$")) {
    # Match "ns" at the end of the string
    as.numeric(str_replace(time_str, "ns", ""))
  } else if (str_detect(time_str, "µs$")) {
    # Match "µs" at the end
    as.numeric(str_replace(time_str, "µs", "")) * 1e3
  } else if (str_detect(time_str, "ms$")) {
    # Match "ms" at the end
    as.numeric(str_replace(time_str, "ms", "")) * 1e6
  } else if (str_detect(time_str, "s$")) {
    # Match "s" at the end
    as.numeric(str_replace(time_str, "s", "")) * 1e9
  } else if (str_detect(time_str, "µ$")) {
    #micro sign
    as.numeric(str_replace(time_str, "µ", "")) * 1e3
  } else {
    NA  # Handle unknown units
  }
}

parseT_into_string <- function(x) {
  ifelse(x >= 1e9, # nanoseconds to seconds
         paste0(round(x / 1e9, 2), " s"), ifelse(
           x >= 1e6,
           # nanoseconds to milliseconds
           paste0(round(x / 1e6, 2), " ms"),
           ifelse(x >= 1e3, # nanoseconds to microseconds
                  paste0(round(x / 1e3, 2), " µs"), paste0(round(x, 0), " ns"))
         )) # nanoseconds
}
# Apply the conversion function
data_numeric <- sapply(data_strings, read_into_ns)

ifelse(test = 0!=length(data_numeric[is.na(data_numeric)]),yes = stop("na's in data"), no = "no numeric problems with the data")
df[is.na(data_numeric),]
df <- data.frame(df, V3 = data_numeric)
colnames(df) <- c("stringlike", "testType", "Time_ns")
count_outliers <- df %>% dplyr:::filter(Time_ns >= 3e+04) %>% length()
df <- df %>% dplyr:::filter(Time_ns <= 3e+04)

# Get the unique levels of testType for the loop
test_levels <- levels(df$testType)

# Define colors for each level (customize as needed)
colors <- c(
  "#CCFFCC",
  "#F9FFC2",
  "#F8DD59",
  "#DEA954",
  "#5E807F",
  "#99621E",
  "#FFAA5C",
  "#FFE3E0",
  "#D4C59D",
  "#A1745E",
  "#664E3B",
  "#34474F",
  "#588157",
  "#90B050",
  "#E67E22",
  "#B92E40",
  "#954FB2",
  "#677D8C",
  "#BDC3C7",
  "#F39C12",
  "#E74C3C"
)
scales:::show_col(sample(colors))
sampled_colors <- sample(colors)[1:length(test_levels)]
color_mapping <- setNames(sampled_colors, test_levels) # Map levels to colors

summarystat<- function(x) {
y <- x %>% is.na()
x = x[!y]
max(x)
   z1 <- mean(x)
  z2 <- median(x)
  z3 <- sd(x)
  z1; z2 ; z3
  return(list(mean=z1, median=z2, sd=z3))
}
summarystat(data_numeric)
#bigger than 10^5
# The stacked histogram plot:
execfrequencies <- ggplot(df, aes(x = Time_ns, fill = testType)) +  # Fill by testType
  geom_histogram(bins = 120,
                 color = "#AAAAAA",
                 linewidth = 0.5)+ #bins=30 to have a default value, you can change it as you want
  scale_y_continuous(name = "Frequency (Total bin count Log)", position = "right",
                     #labels = function(x) {round(log10(x)*(70000/30))},
                     #breaks= c(0,10000,20000,30000,40000,50000,60000,70000),
                     #transform = "log10",
  )+
  scale_x_continuous(
     labels = parseT_into_string,
    name = "Execution time",
    #transform = "log10",
    #minor_breaks = c(50,500,5000,5000),
    #breaks = c(1,10,100,1000,10000,100000,1000000,10000000),
    guide = guide_axis(minor.ticks = T),
  )+
  #annotation_logticks(color=  "#FFAA5C") +
  scale_fill_manual(values = color_mapping, name = "testType") +
  labs(title = "Exec times for Logging Functions in Pydpll/Errorutils",subtitle = sprintf("n= %d",length(data_numeric))) +
  theme_minimal() +
  theme(
    plot.background = element_rect(fill = "#101818"),
    legend.position =  "bottom",
    legend.text = element_text(color = "#D6BA7C"),
    legend.title = element_text(color = "orange"),
    legend.key = element_rect(
      fill = "white",
      linewidth = 3,
      color = "transparent"
    ),
    axis.text = element_text(
      color = "#D6BA7C",
      margin = margin(t = 5, r = 5, unit = "lines")
    ),
    axis.line = element_blank(),
    panel.grid =  element_line(color = "#222222"),
    axis.title = element_text(color = "#D6BA7C", size = "12"),
    plot.title = element_text(
      color = "orange",
      size = 18,
      family = "Calistoga",
      margin = margin(t = 10, b = 15),
      hjust = 0.5
    ),
    plot.subtitle = element_text(color = "orange", size = 11, family = "roboto", hjust= 0.5)
  )

execDistributions <-df %>% 
  ggplot(aes(x=testType, y=Time_ns, fill=testType)) +
  geom_boxplot(color="#AAAAAA") +
  scale_y_continuous(
    name = "Execution time",
    labels = parseT_into_string
  )+   coord_flip() +
  scale_fill_manual(values = color_mapping, name = "testType") +
  labs(title = "Exec times for Logging Functions in Pydpll/Errorutils",subtitle = "n=226546") +
  theme_minimal() +
  theme(
    plot.background = element_rect(fill = "#101818"),
    legend.position =  "hidden",
    legend.text = element_text(color = "#D6BA7C"),
    legend.title = element_text(color = "orange"),
    legend.key = element_rect(
      fill = "white",
      linewidth = 3,
      color = "transparent"
    ),
    axis.text = element_text(
      color = "#D6BA7C",
      margin = margin(t = 5, r = 5, unit = "lines")
    ),
    axis.line = element_blank(),
    panel.grid =  element_line(color = "#222222"),
    axis.title = element_text(color = "#D6BA7C", size = "12"),
    plot.title = element_text(
      color = "orange",
      size = 18,
      family = "Calistoga",
      margin = margin(t = 10, b = 15),
      hjust = 0.5
    ),
    plot.subtitle = element_text(color = "orange", size = 11, family = "roboto", hjust= 0.5)
  )







execfrequencies
execDistributions
ggsave("logrus_distributions.png",execDistributions, scale =1.2 ,dpi = 250,device = "png",width = 1950,height = 1490, units = "px")
ggsave("logrus_frequencies.png",execfrequencies, scale =1.2 ,dpi = 250,device = "png",width = 1950,height = 1490, units = "px")
#  Overlapping hist ----
ggplot(df, aes(x = Time_ns)) +
  scale_x_log10() +  # Log scale for x-axis
  labs(title = "Overlapping Histograms of Time_ns by testType (Log Scale)",
       x = "Time (ns) - Log Scale",
       y = "Frequency") +
  theme_bw() +  # Clean theme
  theme(legend.position = "bottom") + # Legend at the bottom
ggsave()
  # Loop through each level to create a separate histogram layer
  mapply(function(level, color) {
    geom_histogram(data = df[df$testType == level, ], # Filter data for this level
                   aes(fill = level), # Use fill aesthetic for color mapping
                   alpha = 0.2,       # Opacity (0.2 for 20%)
                   color = "black",   # Border color
                   bins = 30,
                   position = "identity") # Important for overlapping
  }, test_levels, color_mapping) + # Use mapply to pass levels and colors
  scale_fill_manual(values = color_mapping, name = "testType")
