#!/bin/bash

# This script contains functions which is added to .bashrc file of the VM for
# additional features for docker caching for warpbuild.

docker() {
    # Check if docker caching should be disabled
    local disable_cache=$(echo "$WARPBUILD_DOCKER_DISABLE_CACHE" | tr '[:upper:]' '[:lower:]')
    if [[ "$disable_cache" == "1" ]] || [[ "$disable_cache" == "true" ]]; then
        # If caching is disabled, just pass all arguments to the original docker command
        command docker "$@"
        return
    fi

    # Define cache locations
    local cache_from=$WARPBUILD_DOCKER_S3_CACHE_FROM
    local cache_to=$WARPBUILD_DOCKER_S3_CACHE_TO
    local builder=$WARPBUILD_DOCKER_BUILDER

    # Check if the first argument is "build" or the first two are "buildx build"
    if [[ "$1 $2" == "buildx build" ]]; then
        local add_cache_from=true
        local add_cache_to=true
        local add_builder=true

        # Check all passed arguments for cache-from and cache-to
        for arg in "$@"; do
            if [[ "$arg" == "--cache-from="* ]]; then
                add_cache_from=false
            fi
            if [[ "$arg" == "--cache-to="* ]]; then
                add_cache_to=false
            fi
            if [[ "$arg" == "--builder="* ]]; then
                add_builder=false
            fi
        done

        # Build the command starting with the original docker command
        local cmd="command docker"

        # Append all original arguments
        for arg in "$@"; do
            cmd+=" $arg"
        done

        # Append cache-from and cache-to if they were not set
        if $add_cache_from; then
            cmd+=" --cache-from=$cache_from"
        fi
        if $add_cache_to; then
            cmd+=" --cache-to=$cache_to"
        fi
        if $add_builder; then
            cmd+=" --builder=$builder"
        fi

        # Execute the constructed command
        eval $cmd
    else
        # If not a build command, just pass all arguments to the original docker command
        command docker "$@"
    fi
}
