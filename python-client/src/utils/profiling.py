from functools import wraps
import time
import psutil
from .logging import logger

def print_memory_stats(prefix: str = "") -> None:
    process = psutil.Process()
    memory_info = process.memory_info()
    logger.info(
        f"{prefix}Memory usage: {memory_info.rss / (1024 * 1024):.2f} MB, "
        f"Virtual Memory: {memory_info.vms / (1024 * 1024):.2f} MB"
    )

def timing_decorator(prefix: str = ""):
    def internal_timing_decorator(f):
        @wraps(f)
        def decorator():
            start = time.perf_counter()
            original_return_val = f()
            end = time.perf_counter()
            logger.info("{} - {}: {:.3f}s".format(prefix, f.__name__, end - start))
            return original_return_val
        return decorator
    return internal_timing_decorator